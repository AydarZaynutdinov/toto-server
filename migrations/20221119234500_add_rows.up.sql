INSERT INTO sku_configs (id, package, country_code, percentile_min, percentile_max, main_sku)
VALUES (gen_random_uuid(), 'com.softinit.iquitos.mainapp', 'US', 0, 25, 'rdm_premium_v3_020_tria l_7d_monthly'),
       (gen_random_uuid(), 'com.softinit.iquitos.mainapp', 'US', 25, 50, 'rdm_premium_v3_030_tria l_7d_monthly'),
       (gen_random_uuid(), 'com.softinit.iquitos.mainapp', 'US', 50, 75, 'rdm_premium_v3_100_tria l_7d_yearly'),
       (gen_random_uuid(), 'com.softinit.iquitos.mainapp', 'US', 75, 100, 'rdm_premium_v3_150_tria l_7d_yearly'),
       (gen_random_uuid(), 'com.softinit.iquitos.mainapp', 'ZZ', 0, 100, 'rdm_premium_v3_050_tria l_7d_yearly');
