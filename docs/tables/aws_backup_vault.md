# Table: aws_backup_vault

In AWS Backup, a backup vault is a container that you organize your backups in. You can use backup vaults to set the AWS Key Management Service (AWS KMS) encryption key that is used to encrypt backups in the backup vault and to control access to the backups in the backup vault. If you require different encryption keys or access policies for different groups of backups, you can optionally create multiple backup vaults. Otherwise, you can have all your backups organized in the default backup vault.

## Examples

### Basic Info

```sql
select
  name,
  arn,
  creation_date
from
  aws_backup_vault;
```


### List of backup_vaults older than 90 days

```sql
select
  name,
  arn,
  creation_date
from
  aws_backup_vault
where
  creation_date <= (current_date - interval '90' day)
  order by
  creation_date;
```



