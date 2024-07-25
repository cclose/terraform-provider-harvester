package image

import (
	harvsterv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/harvester/terraform-provider-harvester/internal/util"
	"github.com/harvester/terraform-provider-harvester/pkg/constants"
)

const (
	// field descriptions
	DisplayNameDescription = "name to display in the Harvester UI"
	ImageSourceTypeDescription = "source of the virtual machine image. Must be one of the following:\n" +
		"- `%s`: Download the image from a remote URL.\n" +
		"- `%s`: Upload the image from a local file.\n" +
		"- `%s`: Export the image from an existing volume. This clones the volume and is useful if you have an existing VM from which you want to make an image. Requires " +
	        "        %s and %s to specify the volume to clone."
	PVCNamespaceDescription = "namespace of the PVC to export as a new image when using `%s=%s`."
	PVCNameDescription = "name of the PVC to export as a new image when using `%s=%s`."
	StorageClassNameDescription = "the storage class to use to store the image on harvester. If not specified, will use the default storage class."
	URLDescription = "supports the `raw` and `qcow2` image formats which are supported by [qemu](https://www.qemu.org/docs/master/system/images.html#disk-image-file-formats). Bootable ISO images can also be used and are treated like `raw` images."
	// Descriptions for read-only fields
	ProgressDescription = "progress percentage of the image while it is being imported."
	SizeDescription = "size of the image"
	StorageClassParametersDescription = "[parameters](https://docs.harvesterhci.io/v1.3/advanced/storageclass#parameters-tab) of the storage class storing the image."
	VolumeStorageClassNameDescription = "name of the [Storage Class](https://docs.harvesterhci.io/v1.3/advanced/storageclass) for the image volume. This is used when mounting this image to new vms."
)

func Schema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		constants.FieldImageDisplayName: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description: DisplayNameDescription,
		},
		constants.FieldImageURL: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Description:  URLDescription,
		},
		constants.FieldImagePVCNamespace: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: util.IsValidName,
			Description: fmt.Sprintf(PVCNamespaceDescription, constants.FieldImageSourceType, harvsterv1.VirtualMachineImageSourceTypeExportVolume)
		},
		constants.FieldImagePVCName: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: util.IsValidName,
			Description: fmt.Sprintf(PVCNameDescription, constants.FieldImageSourceType, harvsterv1.VirtualMachineImageSourceTypeExportVolume)
		},
		constants.FieldImageSourceType: {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				harvsterv1.VirtualMachineImageSourceTypeDownload,
				harvsterv1.VirtualMachineImageSourceTypeUpload,
				harvsterv1.VirtualMachineImageSourceTypeExportVolume,
			}, false),
			Description: fmt.Sprintf(ImageSourceTypeDescription, 
				ImageSourceTypeDownload, 
				ImageSourceTypeUpload, 
				ImageSourceTypeExportVolume,
				constants.FieldImagePVCNamespace, constants.FieldImagePVCName),
		},
		constants.FieldImageProgress: {
			Type:     schema.TypeInt,
			Computed: true,
			Description: ProgressDescription,
		},
		constants.FieldImageSize: {
			Type:     schema.TypeInt,
			Computed: true,
			Description: SizeDescription,
		},
		constants.FieldImageStorageClassName: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: util.IsValidName,
			Description: StorageClassNameDescription,
		},
		constants.FieldImageStorageClassParameters: {
			Type:     schema.TypeMap,
			Computed: true,
			Description: StorageClassParametersDescription,
		},
		constants.FieldImageVolumeStorageClassName: {
			Type:     schema.TypeString,
			Computed: true,
			Description: VolumeStorageClassNameDescription,
		},
	}
	util.NamespacedSchemaWrap(s, false)
	return s
}

func DataSourceSchema() map[string]*schema.Schema {
	s := util.DataSourceSchemaWrap(Schema())
	s[constants.FieldCommonName].Required = false
	s[constants.FieldCommonName].Optional = true
	s[constants.FieldImageDisplayName].Computed = false
	s[constants.FieldImageDisplayName].Optional = true
	return s
}
