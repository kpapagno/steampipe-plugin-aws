# Table: aws_cloudfront_distribution

AWS Systems CloudFront is a web service that speeds up distribution of your static and dynamic web content, such as .html, .css, .js, and image files, to your users.

## Examples

### Basic info

```sql
select
	id,
	enabled,
	e_tag,
	status,
	domain_name
from
	aws_cloudfront_distribution;
```