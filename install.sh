arch=$(uname -m)
case $arch in
x86_64)
  ARCH=amd64
  ;;
i386)
  ARCH=386
  ;;
*)
  echo "no build for this architecture: $arch"
  exit 1
esac

kernel=$(uname -s)
if [ $kernel != "Darwin" ] && [ $kernel != "Linux" ]; then
  echo "no build for this kernel: $kernel"
  exit 1
fi

KERNEL=$(echo $kernel | tr "[:upper:]" "[:lower:]")
terraform_plugins="$HOME/.terraform.d/plugins/${KERNEL}_$ARCH/"

# IFS= preserve newlines
IFS= manifest=$(curl -s https://api.github.com/repos/awsiv/terraform-provider-rdsdataservice/releases/latest)

url=$(echo $manifest \
| grep "browser_download_url.*${KERNEL}_${ARCH}" \
| cut -d '"' -f 4 \
)
version=$(echo $manifest \
| grep tag_name \
| cut -d '"' -f 4 \
)

if [ -z ${url} ]; then
  echo "no build for this kernel/arch: ${KERNEL}_${ARCH}"
  exit 1
fi

dest_file="terraform-provider-rdsdataservice_$version"
origin_file="terraform-provider-rdsdataservice_${version}_${KERNEL}_${ARCH}.tar.gz"
curl $url -L -o $origin_file
tar xvf $origin_file
mv terraform-provider-rdsdataservice $dest_file
chmod +x $dest_file

mkdir -p $terraform_plugins
mv $dest_file $terraform_plugins/

