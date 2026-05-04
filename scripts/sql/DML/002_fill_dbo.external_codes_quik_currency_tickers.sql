MERGE INTO dbo.external_codes AS tgt
USING (
    SELECT ext_system_id, ext_code, ext_code_type_id, internal_id
    FROM (VALUES
        (1, N'GLD', 1, 959),   -- XAU
        (1, N'SLV', 1, 961),   -- XAG
        (1, N'PLT', 1, 962),   -- XPT
        (1, N'PLD', 1, 964)    -- XPD
    ) AS t(ext_system_id, ext_code, ext_code_type_id, internal_id)
) AS src ON tgt.ext_system_id = src.ext_system_id
    AND tgt.ext_code = src.ext_code
    AND tgt.ext_code_type_id = src.ext_code_type_id
WHEN NOT MATCHED BY TARGET THEN
    INSERT (ext_system_id, ext_code, ext_code_type_id, internal_id)
    VALUES (src.ext_system_id, src.ext_code, src.ext_code_type_id, src.internal_id);
GO
