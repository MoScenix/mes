-- Glucose meter MES demo data. Destructive: preserves only the root administrator.
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE inventory_flow_item_units;
TRUNCATE TABLE inventory_flow_items;
TRUNCATE TABLE inventory_flows;
TRUNCATE TABLE item_units;
TRUNCATE TABLE engineering_orders;
TRUNCATE TABLE process_items;
TRUNCATE TABLE processes;
TRUNCATE TABLE items;
TRUNCATE TABLE work_order;
TRUNCATE TABLE messages;
TRUNCATE TABLE history;
DELETE FROM `user` WHERE user_account <> 'root';

SET FOREIGN_KEY_CHECKS = 1;

SET @root_id = (SELECT id FROM `user` WHERE user_account = 'root' LIMIT 1);
SET @password_hash = (SELECT password_hash FROM `user` WHERE id = @root_id);

INSERT INTO `user`(created_at, updated_at, name, password_hash, user_account, user_avatar, user_profile, user_role) VALUES
(NOW() - INTERVAL 90 DAY, NOW(), '采购专员-林悦', @password_hash, 'purchase_demo', '', '血糖仪原料采购', 'purchase'),
(NOW() - INTERVAL 90 DAY, NOW(), '生产组长-周工', @password_hash, 'leader_demo', '', '血糖仪装配线负责人', 'leader'),
(NOW() - INTERVAL 90 DAY, NOW(), '工艺工程师-陈工', @password_hash, 'process_demo', '', '血糖仪生产工艺', 'process_engineer'),
(NOW() - INTERVAL 90 DAY, NOW(), '仓库管理员-吴姐', @password_hash, 'warehouse_demo', '', '原料与成品仓', 'warehouse_admin'),
(NOW() - INTERVAL 90 DAY, NOW(), '质检员-李师傅', @password_hash, 'worker_demo', '', '终检与入库扫码', 'worker'),
(NOW() - INTERVAL 90 DAY, NOW(), '销售专员-王敏', @password_hash, 'sales_demo', '', '成品出库协调', 'sales');

SET @purchase = (SELECT id FROM `user` WHERE user_account='purchase_demo');
SET @leader = (SELECT id FROM `user` WHERE user_account='leader_demo');
SET @process = (SELECT id FROM `user` WHERE user_account='process_demo');
SET @warehouse = (SELECT id FROM `user` WHERE user_account='warehouse_demo');
SET @worker = (SELECT id FROM `user` WHERE user_account='worker_demo');
SET @sales = (SELECT id FROM `user` WHERE user_account='sales_demo');

INSERT INTO items(created_at,updated_at,name,unit,description,total_count,in_stock_count,reserved_count,out_stock_count,pending_count,qualified_count,unqualified_count,available_count) VALUES
(NOW()-INTERVAL 120 DAY,NOW(),'GM-100 智能血糖仪','台','家用智能血糖监测终端',170,150,0,20,12,150,8,132),
(NOW()-INTERVAL 120 DAY,NOW(),'BG-PCB-01 血糖仪主控板','片','低功耗主控与信号采集板',280,220,20,40,0,280,0,200),
(NOW()-INTERVAL 120 DAY,NOW(),'BG-SEN-02 电化学传感模块','套','血糖试纸电信号采集模块',260,205,15,40,0,258,2,188),
(NOW()-INTERVAL 120 DAY,NOW(),'LCD-128 液晶显示屏','块','128x64 医疗级液晶屏',300,245,15,40,0,296,4,226),
(NOW()-INTERVAL 120 DAY,NOW(),'CASE-GM100 ABS 外壳','套','医用级阻燃 ABS 上下壳',320,260,20,40,0,316,4,240),
(NOW()-INTERVAL 120 DAY,NOW(),'BAT-CR2032 纽扣电池','枚','CR2032 3V 纽扣电池',400,335,25,40,0,400,0,310),
(NOW()-INTERVAL 120 DAY,NOW(),'STRIP-G50 血糖试纸','盒','50片装配套血糖试纸',240,210,10,20,0,238,2,198);

SET @meter = (SELECT id FROM items WHERE name LIKE 'GM-100%');
SET @pcb = (SELECT id FROM items WHERE name LIKE 'BG-PCB%');
SET @sensor = (SELECT id FROM items WHERE name LIKE 'BG-SEN%');
SET @lcd = (SELECT id FROM items WHERE name LIKE 'LCD-128%');
SET @case = (SELECT id FROM items WHERE name LIKE 'CASE-GM100%');
SET @battery = (SELECT id FROM items WHERE name LIKE 'BAT-CR2032%');
SET @strip = (SELECT id FROM items WHERE name LIKE 'STRIP-G50%');

INSERT INTO processes(created_at,updated_at,item_id,owner_user_id,name,description,status) VALUES
(NOW()-INTERVAL 75 DAY,NOW()-INTERVAL 20 DAY,@meter,@process,'GM-100 标准装配与校准工艺','主板烧录、传感模块焊接、整机装配、三点校准、终检',2),
(NOW()-INTERVAL 45 DAY,NOW()-INTERVAL 8 DAY,@meter,@process,'GM-100 快速装配工艺','演示线优化工艺，增加自动校准与扫码追溯',2);
SET @process_std = (SELECT id FROM processes ORDER BY id LIMIT 1);
SET @process_fast = (SELECT id FROM processes ORDER BY id DESC LIMIT 1);

INSERT INTO process_items(created_at,updated_at,process_id,consume_item_id,quantity) VALUES
(NOW()-INTERVAL 75 DAY,NOW(),@process_std,@pcb,1),(NOW()-INTERVAL 75 DAY,NOW(),@process_std,@sensor,1),
(NOW()-INTERVAL 75 DAY,NOW(),@process_std,@lcd,1),(NOW()-INTERVAL 75 DAY,NOW(),@process_std,@case,1),
(NOW()-INTERVAL 75 DAY,NOW(),@process_std,@battery,1),(NOW()-INTERVAL 45 DAY,NOW(),@process_fast,@pcb,1),
(NOW()-INTERVAL 45 DAY,NOW(),@process_fast,@sensor,1),(NOW()-INTERVAL 45 DAY,NOW(),@process_fast,@lcd,1),
(NOW()-INTERVAL 45 DAY,NOW(),@process_fast,@case,1),(NOW()-INTERVAL 45 DAY,NOW(),@process_fast,@battery,1);

INSERT INTO engineering_orders(created_at,updated_at,leader_user_id,item_id,expected_quantity,qualified_quantity,produced_quantity,description,process_id,unqualified_quantity,status,name) VALUES
(NOW()-INTERVAL 18 DAY,NOW()-INTERVAL 6 DAY,@leader,@meter,60,60,60,'华东渠道首批交付',@process_std,0,2,'PP-202607-A01 华东首批生产计划'),
(NOW()-INTERVAL 12 DAY,NOW()-INTERVAL 3 DAY,@leader,@meter,50,50,52,'连锁药房补货批次',@process_fast,2,2,'PP-202607-A02 药房补货生产计划'),
(NOW()-INTERVAL 7 DAY,NOW()-INTERVAL 1 HOUR,@leader,@meter,80,35,45,'电商渠道本周生产任务',@process_fast,4,2,'PP-202607-B01 电商渠道生产计划'),
(NOW()-INTERVAL 3 DAY,NOW()-INTERVAL 20 MINUTE,@leader,@meter,100,5,13,'社区医疗试点批次',@process_std,2,2,'PP-202607-B02 社区医疗生产计划'),
(NOW()-INTERVAL 1 DAY,NOW(),@leader,@meter,120,0,0,'下周出口预备批次',@process_fast,0,2,'PP-202607-C01 出口预备生产计划'),
(NOW(),NOW(),@leader,@meter,40,0,0,'演示用未提交计划草稿',@process_std,0,1,'PP-DRAFT 体验装生产计划');

SET @order1=(SELECT id FROM engineering_orders WHERE name LIKE 'PP-202607-A01%');
SET @order2=(SELECT id FROM engineering_orders WHERE name LIKE 'PP-202607-A02%');
SET @order3=(SELECT id FROM engineering_orders WHERE name LIKE 'PP-202607-B01%');
SET @order4=(SELECT id FROM engineering_orders WHERE name LIKE 'PP-202607-B02%');

-- 170 produced meters spread across seven days for the dashboard trend.
INSERT INTO item_units(created_at,updated_at,item_id,engineering_order_id,stock_status,quality_status,description)
SELECT
  TIMESTAMP(CURRENT_DATE - INTERVAL (6 - x.day_no) DAY) + INTERVAL (8 + MOD(x.n,9)) HOUR + INTERVAL MOD(x.n*7,60) MINUTE,
  NOW(), @meter,
  CASE WHEN x.seq<=60 THEN @order1 WHEN x.seq<=112 THEN @order2 WHEN x.seq<=157 THEN @order3 ELSE @order4 END,
  CASE WHEN x.seq<=150 THEN 1 ELSE 3 END,
  CASE WHEN x.seq<=150 THEN 2 WHEN x.seq<=162 THEN 1 ELSE 3 END,
  CONCAT('GM100-',DATE_FORMAT(CURRENT_DATE,'%Y%m'),'-',LPAD(x.seq,4,'0'))
FROM (
  SELECT d.day_no, n.n, (d.offset_count+n.n) seq
  FROM (
    SELECT 0 day_no,0 offset_count,12 day_count UNION ALL SELECT 1,12,18 UNION ALL
    SELECT 2,30,25 UNION ALL SELECT 3,55,22 UNION ALL SELECT 4,77,30 UNION ALL
    SELECT 5,107,35 UNION ALL SELECT 6,142,28
  ) d
  JOIN (
    SELECT ones.i + tens.i*10 + 1 n
    FROM (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) ones
    CROSS JOIN (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3) tens
  ) n ON n.n<=d.day_count
) x;

-- Representative raw-material units; item counters above retain the full warehouse totals.
INSERT INTO item_units(created_at,updated_at,item_id,engineering_order_id,stock_status,quality_status,description)
SELECT NOW()-INTERVAL MOD(seq,30) DAY,NOW(),material_id,NULL,stock_status,quality_status,CONCAT(code,'-',LPAD(seq,4,'0'))
FROM (
  SELECT m.material_id,m.code,n.seq,
    CASE WHEN n.seq<=30 THEN 1 ELSE 3 END stock_status,
    CASE WHEN n.seq=40 THEN 3 ELSE 2 END quality_status
  FROM (SELECT @pcb material_id,'PCB' code UNION ALL SELECT @sensor,'SEN' UNION ALL SELECT @lcd,'LCD' UNION ALL SELECT @case,'CASE' UNION ALL SELECT @battery,'BAT' UNION ALL SELECT @strip,'STRIP') m
  CROSS JOIN (SELECT ones.i+tens.i*10+1 seq FROM (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) ones CROSS JOIN (SELECT 0 i UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3) tens) n
) raw_units;

INSERT INTO inventory_flows(created_at,updated_at,from_user_id,to_user_id,flow_type,business_type,flow_status,description,approved_by,approved_at,name) VALUES
(NOW()-INTERVAL 14 DAY,NOW()-INTERVAL 13 DAY,@purchase,@warehouse,1,1,3,'首批电子元器件采购到货',@warehouse,NOW()-INTERVAL 13 DAY,'PO-IN-071 电子元器件采购入库'),
(NOW()-INTERVAL 8 DAY,NOW()-INTERVAL 7 DAY,@purchase,@warehouse,1,1,3,'外壳与包装辅料到货',@warehouse,NOW()-INTERVAL 7 DAY,'PO-IN-086 结构件采购入库'),
(NOW()-INTERVAL 1 DAY,NOW()-INTERVAL 1 HOUR,@purchase,@warehouse,1,1,2,'传感模块紧急补货待审批',0,NULL,'PO-IN-102 传感模块采购入库'),
(NOW(),NOW(),@purchase,@warehouse,1,1,1,'下周电池采购草稿',0,NULL,'PO-DRAFT 电池采购入库'),
(NOW()-INTERVAL 10 DAY,NOW()-INTERVAL 9 DAY,@leader,@warehouse,2,2,3,'A01计划生产领料',@warehouse,NOW()-INTERVAL 9 DAY,'MR-071 A01生产领料'),
(NOW()-INTERVAL 5 DAY,NOW()-INTERVAL 4 DAY,@leader,@warehouse,2,2,3,'B01计划生产领料',@warehouse,NOW()-INTERVAL 4 DAY,'MR-088 B01生产领料'),
(NOW()-INTERVAL 2 HOUR,NOW()-INTERVAL 30 MINUTE,@leader,@warehouse,2,2,2,'社区医疗批次追加领料',0,NULL,'MR-105 B02追加领料'),
(NOW(),NOW(),@leader,@warehouse,2,2,1,'出口批次预领料草稿',0,NULL,'MR-DRAFT C01预领料'),
(NOW()-INTERVAL 6 DAY,NOW()-INTERVAL 6 DAY,@leader,@warehouse,1,3,3,'A01完工合格品入库',@warehouse,NOW()-INTERVAL 6 DAY,'FG-IN-061 A01成品入库'),
(NOW()-INTERVAL 3 DAY,NOW()-INTERVAL 3 DAY,@leader,@warehouse,1,3,3,'A02完工合格品入库',@warehouse,NOW()-INTERVAL 3 DAY,'FG-IN-079 A02成品入库'),
(NOW()-INTERVAL 3 HOUR,NOW()-INTERVAL 20 MINUTE,@leader,@warehouse,1,3,2,'B01阶段性完工入库',0,NULL,'FG-IN-108 B01成品入库'),
(NOW(),NOW(),@leader,@warehouse,1,3,1,'B02阶段性入库草稿',0,NULL,'FG-DRAFT B02成品入库');

INSERT INTO inventory_flow_items(created_at,updated_at,inventory_flow_id,item_id,apply_quantity,finished_quantity)
SELECT f.created_at,f.updated_at,f.id,
  CASE f.business_type WHEN 3 THEN @meter WHEN 2 THEN @pcb ELSE @sensor END,
  CASE WHEN f.business_type=3 THEN 50 WHEN f.business_type=2 THEN 80 ELSE 100 END,
  CASE WHEN f.flow_status=3 THEN CASE WHEN f.business_type=3 THEN 50 WHEN f.business_type=2 THEN 80 ELSE 100 END ELSE 0 END
FROM inventory_flows f;

-- All demo operations are owned by root so every administrator view is populated.
UPDATE inventory_flows
SET from_user_id=@root_id,
    approved_by=CASE WHEN approved_by<>0 THEN @root_id ELSE 0 END;
UPDATE engineering_orders SET leader_user_id=@root_id;
UPDATE processes SET owner_user_id=@root_id;

INSERT INTO work_order(created_at,updated_at,from_user_id,to_user_id,description,status,read_status,name) VALUES
(NOW()-INTERVAL 9 DAY,NOW()-INTERVAL 8 DAY,@leader,@worker,'完成A01批次整机终检与扫码入库',2,2,'任务：A01批次质量终检'),
(NOW()-INTERVAL 4 DAY,NOW()-INTERVAL 3 DAY,@leader,@worker,'优先检验B01批次传感精度与显示功能',2,2,'任务：B01批次抽检'),
(NOW()-INTERVAL 1 DAY,NOW()-INTERVAL 2 HOUR,@leader,@worker,'处理当前12台待检血糖仪',2,1,'任务：清理待检队列'),
(NOW()-INTERVAL 2 DAY,NOW()-INTERVAL 1 DAY,@warehouse,@purchase,'传感模块库存低于安全线，请安排补货',2,1,'提醒：传感模块补货'),
(NOW(),NOW(),@leader,@worker,'演示用质检任务草稿',1,1,'草稿：C01首件检验');

UPDATE work_order SET from_user_id=@root_id;

-- Keep aggregate item metadata aligned with the intended demo presentation.
UPDATE items SET updated_at=NOW() WHERE deleted_at IS NULL;

SELECT 'demo data ready' result,
       (SELECT COUNT(*) FROM `user` WHERE deleted_at IS NULL) users,
       (SELECT COUNT(*) FROM items) items,
       (SELECT COUNT(*) FROM item_units WHERE engineering_order_id IS NOT NULL) produced_units,
       (SELECT COUNT(*) FROM engineering_orders) production_plans,
       (SELECT COUNT(*) FROM inventory_flows) business_flows;
