-- 为 email_audit_applications 表添加外键约束
-- 如果表已经存在但没有外键约束，执行此SQL来添加外键

-- 注意：执行前请确保：
-- 1. email_audit_applications 表已存在
-- 2. sys_users 表已存在
-- 3. 所有 applicant_id 和 auditor_id 的值都存在于 sys_users.id 中（否则会报错）

-- 检查并添加 applicant_id 外键约束
-- 如果外键已存在，会报错，可以忽略或先删除再添加
ALTER TABLE `email_audit_applications`
ADD CONSTRAINT `fk_email_audit_applicant` 
FOREIGN KEY (`applicant_id`) 
REFERENCES `sys_users` (`id`) 
ON DELETE RESTRICT 
ON UPDATE CASCADE;

-- 检查并添加 auditor_id 外键约束
ALTER TABLE `email_audit_applications`
ADD CONSTRAINT `fk_email_audit_auditor` 
FOREIGN KEY (`auditor_id`) 
REFERENCES `sys_users` (`id`) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- 如果外键已存在，需要先删除再添加，可以使用以下语句：
-- ALTER TABLE `email_audit_applications` DROP FOREIGN KEY `fk_email_audit_applicant`;
-- ALTER TABLE `email_audit_applications` DROP FOREIGN KEY `fk_email_audit_auditor`;
-- 然后重新执行上面的 ADD CONSTRAINT 语句

