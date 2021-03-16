package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type distributionInfo = struct {
	cloudfront.DistributionSummary
	DistributionConfig            *cloudfront.DistributionConfig
	ActiveTrustedKeyGroups        *cloudfront.ActiveTrustedKeyGroups
	ActiveTrustedSigners          *cloudfront.ActiveTrustedSigners
	InProgressInvalidationBatches *int64
	Etag                          *string
}

//// TABLE DEFINITION

func tableAwsCloudfrontDistribution(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_distribution",
		Description: "AWS Cloudfront Distribution",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"NoSuchDistribution"}),
			Hydrate:           getCloudfrontDistribution,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsCloudfrontDistribution,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The identifier for the Distribution.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled",
				Description: "Whether the Distribution is enabled to accept user requests for content.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "e_tag",
				Description: "The current version of the distribution's information.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudfrontDistribution,
				Transform:   transform.FromField("Etag"),
			},
			{
				Name:        "status",
				Description: "The current status of the Distribution. When the status is Deployed, the distribution's information is propagated to all CloudFront edge locations.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_time",
				Description: "The date and time the Distribution was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "domain_name",
				Description: "The domain name that corresponds to the Distribution, for example, d111111abcdef8.cloudfront.net.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Maintenance Window",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudfrontDistributionTags,
				Transform:   transform.FromField("Tags.Items"),
			},
			{
				Name:        "comment",
				Description: "The comment originally specified when this Distribution was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "http_version",
				Description: "Specify the maximum HTTP version that you want viewers to use to communicate with CloudFront. The default value for new web Distributions is http2. Viewers that don't support HTTP/2 will automatically use an earlier version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_ipv6_enabled",
				Description: "Whether CloudFront responds to IPv6 DNS requests with an IPv6 address for your Distribution.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsIPV6Enabled"),
			},
			{
				Name:        "alias_icp_recordals",
				Description: "AWS services in China customers must file for an Internet Content Provider (ICP) recordal if they want to serve content publicly on an alternate domain name, also known as a CNAME, that they've added to CloudFront. AliasICPRecordal provides the ICP recordal status for CNAMEs associated with distributions.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AliasICPRecordals"),
			},
			{
				Name:        "custom_error_responses",
				Description: "A complex type that contains zero or more CustomErrorResponses elements.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "default_cache_behavior",
				Description: "A complex type that describes the default cache behavior if you don't specify a CacheBehavior element or if files don't match any of the values of PathPattern in CacheBehavior elements. You must create exactly one default cache behavior.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "origin_groups",
				Description: "A complex type that contains information about origin groups for this distribution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "restrictions",
				Description: "A complex type that identifies ways in which you want to restrict distribution of your content.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "viewer_certificate",
				Description: "A complex type that determines the distribution’s SSL/TLS configuration for communicating with viewers.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "active_trusted_key_groups",
				Description: "A list of key groups, including the identifiers of the public keys in each key group that CloudFront can use to verify the signatures of signed URLs and signed cookies.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudfrontDistribution,
			},
			{
				Name:        "active_trusted_signers",
				Description: "A list of AWS accounts and the identifiers of active CloudFront key pairs in each account that CloudFront can use to verify the signatures of signed URLs and signed cookies.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudfrontDistribution,
			},
			{
				Name:        "price_class",
				Description: "A complex type that contains information about price class for this streaming Distribution.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "web_acl_id",
				Description: "The Web ACL Id (if any) associated with the distribution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WebACLId"),
			},
			{
				Name:        "default_root_object",
				Description: "The object that you want CloudFront to request from your origin.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudfrontDistribution,
				Transform:   transform.FromField("DistributionConfig.DefaultRootObject"),
			},
			{
				Name:        "aliases",
				Description: "The number of CNAME aliases, if any, that you want to associate with this Distribution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cache_behaviors",
				Description: "The number of cache behaviors for this Distribution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "origins",
				Description: "The number of origins for this distribution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "in_progress_invalidation_batches",
				Description: "A list of origins. ",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudfrontDistribution,
				Transform:   transform.FromField("InProgressInvalidationBatches"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudfrontDistributionTags,
				Transform:   transform.FromField("Tags.Items").Transform(cloudfrontDistributionTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsCloudfrontDistribution(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsCloudfrontDistribution")

	// Create session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListDistributionsPages(
		&cloudfront.ListDistributionsInput{},
		func(page *cloudfront.ListDistributionsOutput, isLast bool) bool {
			for _, parameter := range page.DistributionList.Items {
				d.StreamListItem(ctx, distributionInfo{DistributionSummary: *parameter})
			}
			return !isLast
		},
	)

	return nil, err
}

func getCloudfrontDistribution(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudfrontDistribution")

	// Create session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}

	var cloudfrontID string
	if h.Item != nil {
		cloudfrontID = *h.Item.(distributionInfo).Id
	} else {
		cloudfrontID = d.KeyColumnQuals["id"].GetStringValue()
	}

	params := &cloudfront.GetDistributionInput{
		Id: &cloudfrontID,
	}

	op, err := svc.GetDistribution(params)
	if err != nil {
		return nil, err
	}

	return distributionInfo{
		DistributionSummary: cloudfront.DistributionSummary{
			ARN:                  op.Distribution.ARN,
			Id:                   op.Distribution.Id,
			DomainName:           op.Distribution.DomainName,
			Status:               op.Distribution.Status,
			LastModifiedTime:     op.Distribution.LastModifiedTime,
			AliasICPRecordals:    op.Distribution.AliasICPRecordals,
			Enabled:              op.Distribution.DistributionConfig.Enabled,
			HttpVersion:          op.Distribution.DistributionConfig.HttpVersion,
			IsIPV6Enabled:        op.Distribution.DistributionConfig.IsIPV6Enabled,
			PriceClass:           op.Distribution.DistributionConfig.PriceClass,
			WebACLId:             op.Distribution.DistributionConfig.WebACLId,
			Comment:              op.Distribution.DistributionConfig.Comment,
			Origins:              op.Distribution.DistributionConfig.Origins,
			Aliases:              op.Distribution.DistributionConfig.Aliases,
			CacheBehaviors:       op.Distribution.DistributionConfig.CacheBehaviors,
			DefaultCacheBehavior: op.Distribution.DistributionConfig.DefaultCacheBehavior,
		},
		DistributionConfig:            op.Distribution.DistributionConfig,
		ActiveTrustedKeyGroups:        op.Distribution.ActiveTrustedKeyGroups,
		ActiveTrustedSigners:          op.Distribution.ActiveTrustedSigners,
		InProgressInvalidationBatches: op.Distribution.InProgressInvalidationBatches,
		Etag:                          op.ETag,
	}, nil
}

func getCloudfrontDistributionTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudfrontDistributionTags")

	// Create session
	svc, err := CloudFrontService(ctx, d)
	if err != nil {
		return nil, err
	}
	data := h.Item.(distributionInfo)

	// Build the params
	params := &cloudfront.ListTagsForResourceInput{
		Resource: data.ARN,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func cloudfrontDistributionTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("cloudfrontDistributionTagListToTurbotTags")
	tagList := d.Value.([]*cloudfront.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	} else {
		return nil, nil
	}

	return turbotTagsMap, nil
}
