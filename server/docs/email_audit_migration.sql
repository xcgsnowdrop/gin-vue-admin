-- 邮件审核申请表创建SQL
-- 执行此SQL创建邮件审核申请表

CREATE TABLE IF NOT EXISTS `email_audit_applications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `applicant_id` bigint unsigned NOT NULL COMMENT '申请人ID（关联sys_users.id）',
  `applicant_time` datetime NOT NULL COMMENT '申请时间',
  `auditor_id` bigint unsigned DEFAULT NULL COMMENT '审核人ID（关联sys_users.id）',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `audit_comment` text COMMENT '审核说明',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1-待审核，2-通过，3-拒绝，4-待修改',
  `email_type` tinyint NOT NULL COMMENT '邮件类型',
  `email_title` json NOT NULL COMMENT '邮件标题（多语言JSON）',
  `email_content` json NOT NULL COMMENT '邮件内容（多语言JSON）',
  `email_attachments` json DEFAULT NULL COMMENT '邮件附件（JSON数组）',
  `email_remark` varchar(500) DEFAULT NULL COMMENT '邮件备注',
  `start_time` bigint DEFAULT NULL COMMENT '开始生效时间（时间戳秒）',
  `area_ids` varchar(500) DEFAULT NULL COMMENT '生效区服列表（逗号分隔）',
  `max_reg_time` bigint DEFAULT NULL COMMENT '最大注册时间（时间戳秒）',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_applicant_id` (`applicant_id`),
  KEY `idx_auditor_id` (`auditor_id`),
  KEY `idx_status` (`status`),
  KEY `idx_applicant_time` (`applicant_time`),
  KEY `idx_deleted_at` (`deleted_at`),
  -- 外键约束：applicant_id 引用 sys_users.id
  CONSTRAINT `fk_email_audit_applicant` FOREIGN KEY (`applicant_id`) 
    REFERENCES `sys_users` (`id`) 
    ON DELETE RESTRICT 
    ON UPDATE CASCADE,
  -- 外键约束：auditor_id 引用 sys_users.id
  CONSTRAINT `fk_email_audit_auditor` FOREIGN KEY (`auditor_id`) 
    REFERENCES `sys_users` (`id`) 
    ON DELETE SET NULL 
    ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邮件审核申请表';

