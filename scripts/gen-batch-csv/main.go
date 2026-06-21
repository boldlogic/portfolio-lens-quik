package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type templateRow struct {
	fields     []string
	clientIdx  int
	limitType  string
	ticker     string
	settleCode string
	firmCode   string
	balance    string
	isin       string
	position   string
	tradeAcct  string
	acqCcy     string
}

func main() {
	templatePath := flag.String("template", "docs/батч.csv", "путь к CSV-шаблону")
	outputPath := flag.String("output", "docs/батч-20mb.csv", "путь к выходному CSV")
	targetMB := flag.Float64("target-mb", 20, "целевой размер файла в МБ")
	clientStart := flag.Int("client-start", 1_000_000, "начальный числовой client_code")
	flag.Parse()

	targetBytes := int64(*targetMB * 1024 * 1024)

	rows, header, err := loadTemplate(*templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения шаблона: %v\n", err)
		os.Exit(1)
	}

	out, err := os.Create(*outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка создания файла: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	w := csv.NewWriter(out)
	w.Comma = ';'

	if err := w.Write(header); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка записи заголовка: %v\n", err)
		os.Exit(1)
	}

	var (
		dataRows  uint64
		clients   int
		written   int64
		headerLen = int64(len(strings.Join(header, ";")) + 1)
	)

	for clientN := *clientStart; written < targetBytes; clientN++ {
		clientCode := strconv.Itoa(clientN)
		clients++

		for _, tpl := range rows {
			record := make([]string, len(tpl.fields))
			copy(record, tpl.fields)
			record[tpl.clientIdx] = clientCode

			if err := validateRow(tpl, clientCode); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка валидации client_code=%s: %v\n", clientCode, err)
				os.Exit(1)
			}

			if err := w.Write(record); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка записи строки: %v\n", err)
				os.Exit(1)
			}
			dataRows++
		}

		w.Flush()
		if err := w.Error(); err != nil {
			fmt.Fprintf(os.Stderr, "ошибка flush: %v\n", err)
			os.Exit(1)
		}

		info, err := out.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка stat: %v\n", err)
			os.Exit(1)
		}
		written = info.Size()
		if written < headerLen {
			written = headerLen
		}
	}

	fmt.Printf("файл: %s\n", *outputPath)
	fmt.Printf("размер: %d байт (цель: %d)\n", written, targetBytes)
	fmt.Printf("строк данных: %d\n", dataRows)
	fmt.Printf("клиентов: %d\n", clients)
	fmt.Printf("строк на клиента: %d\n", len(rows))
}

func loadTemplate(path string) ([]templateRow, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	if len(records) < 2 {
		return nil, nil, fmt.Errorf("шаблон пуст или без строк данных")
	}

	header := records[0]
	colIdx := columnIndex(header)
	clientIdx, ok := colIdx["client_code"]
	if !ok {
		return nil, nil, fmt.Errorf("в шаблоне нет колонки client_code")
	}

	seen := make(map[[32]byte]int)
	rows := make([]templateRow, 0, len(records)-1)

	for i, record := range records[1:] {
		lineNum := i + 2
		tpl, err := parseTemplateRow(header, colIdx, clientIdx, record, lineNum)
		if err != nil {
			return nil, nil, err
		}

		limit, err := newLimitFromTemplate(tpl, tpl.fields[clientIdx])
		if err != nil {
			return nil, nil, fmt.Errorf("строка %d: %w", lineNum, err)
		}

		h := limit.KeyHash()
		if prev, dup := seen[h]; dup {
			return nil, nil, fmt.Errorf("дубликат в шаблоне: строки %d и %d", prev, lineNum)
		}
		seen[h] = lineNum

		rows = append(rows, tpl)
	}

	rows, err = appendOtcRows(rows, colIdx, clientIdx, seen)
	if err != nil {
		return nil, nil, err
	}

	return rows, header, nil
}

func appendOtcRows(rows []templateRow, colIdx map[string]int, clientIdx int, seen map[[32]byte]int) ([]templateRow, error) {
	typeIdx, ok := colIdx["limit_type"]
	if !ok {
		return nil, fmt.Errorf("в шаблоне нет колонки limit_type")
	}
	tradeIdx, ok := colIdx["trade_account"]
	if !ok {
		return nil, fmt.Errorf("в шаблоне нет колонки trade_account")
	}

	otcRows := make([]templateRow, 0)
	for _, tpl := range rows {
		if tpl.limitType != string(quik.LimitTypeSecurities) {
			continue
		}

		derived := tpl
		derived.fields = make([]string, len(tpl.fields))
		copy(derived.fields, tpl.fields)
		derived.fields[typeIdx] = string(quik.LimitTypeSecuritiesOtc)
		derived.fields[tradeIdx] = "OTC"
		derived.limitType = string(quik.LimitTypeSecuritiesOtc)
		derived.tradeAcct = "OTC"

		limit, err := newLimitFromTemplate(derived, derived.fields[clientIdx])
		if err != nil {
			return nil, fmt.Errorf("security_otc из %s/%s: %w", tpl.ticker, tpl.settleCode, err)
		}

		h := limit.KeyHash()
		if prev, dup := seen[h]; dup {
			return nil, fmt.Errorf("дубликат security_otc: конфликт со строкой шаблона %d", prev)
		}
		seen[h] = 0

		otcRows = append(otcRows, derived)
	}

	return append(rows, otcRows...), nil
}

func columnIndex(header []string) map[string]int {
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[h] = i
	}
	return idx
}

func getField(colIdx map[string]int, record []string, name string) string {
	i, ok := colIdx[name]
	if !ok || i >= len(record) {
		return ""
	}
	return record[i]
}

func parseTemplateRow(header []string, colIdx map[string]int, clientIdx int, record []string, lineNum int) (templateRow, error) {
	if len(record) != len(header) {
		return templateRow{}, fmt.Errorf("строка %d: ожидалось %d полей, получено %d", lineNum, len(header), len(record))
	}

	fields := make([]string, len(record))
	copy(fields, record)

	return templateRow{
		fields:     fields,
		clientIdx:  clientIdx,
		limitType:  getField(colIdx, record, "limit_type"),
		ticker:     getField(colIdx, record, "ticker"),
		settleCode: getField(colIdx, record, "settle_code"),
		firmCode:   getField(colIdx, record, "firm_code"),
		balance:    getField(colIdx, record, "balance"),
		isin:       getField(colIdx, record, "isin"),
		position:   getField(colIdx, record, "position_code"),
		tradeAcct:  getField(colIdx, record, "trade_account"),
		acqCcy:     getField(colIdx, record, "acquisition_currency"),
	}, nil
}

func ptrIfNonEmpty(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	v := s
	return &v
}

func newLimitFromTemplate(tpl templateRow, clientCode string) (quik.Limit, error) {
	balance, err := decimal.NewFromString(tpl.balance)
	if err != nil {
		return quik.Limit{}, fmt.Errorf("balance: %w", err)
	}

	return quik.NewLimit(
		tpl.limitType,
		clientCode,
		tpl.ticker,
		ptrIfNonEmpty(tpl.position),
		tpl.settleCode,
		ptrIfNonEmpty(tpl.tradeAcct),
		tpl.firmCode,
		balance,
		ptrIfNonEmpty(tpl.acqCcy),
		ptrIfNonEmpty(tpl.isin),
	)
}

func validateRow(tpl templateRow, clientCode string) error {
	_, err := newLimitFromTemplate(tpl, clientCode)
	return err
}
