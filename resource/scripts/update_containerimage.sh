cd container/ &&  make docker && cd ..

docker save aruba-iotops-example-ble:1.0.0-release > ${imagename}.tar

eval "curl '${url}/iot_operations/api/v1/adp/images/maxversion?pageNumber=1&pageSize=1000' \
    -H 'authorization: Bearer ${token}' -o maxversion.json"

cat maxversion.json

md5val=$(md5sum ${imagename}.tar | awk '{print $1}')
version=$(jq '.content | .[] | select (.name == "'${imagename}'")'.version+1 maxversion.json)

if [ -z ${version} ]; then
    version=1
fi

data='{ "img_name": "'${imagename}'", "img_type": "RUNS_ON_COLLECTOR", "md5": "'${md5val}'", "version": '${version}' }'
echo $data

eval "curl '${url}/iot_operations/api/v1/adp/images' \
  -X 'POST' \
  -H 'authorization: Bearer ${token}' \
  -H 'content-type: application/json' \
  -H 'referer: ${url}/swagger/central/' \
  --data-raw '${data}' \
  -o imageupload.json"

cat imageupload.json

echo "curl --location $(jq .post_destination imageupload.json) --header 'Content-Type: multipart/form-data' --form 'key=\"$(jq -r --arg prefix "${imagename}.tar" '.key + $prefix' imageupload.json)\"' --form 'success_action_status="201"' --form 'Content-Type=\"application/x-tar\"' --form 'x-amz-meta-uuid="$(jq '."x-amz-meta-uuid"' imageupload.json)"' --form 'x-amz-credential="$(jq '."x-amz-credential"' imageupload.json)"' --form 'x-amz-algorithm="$(jq '."x-amz-algorithm"' imageupload.json)"' --form 'x-amz-date="$(jq '."x-amz-date"' imageupload.json)"' --form 'x-amz-meta-tag="$(jq '."x-amz-meta-tag"' imageupload.json)"' --form 'policy="$(jq .policy imageupload.json)"' --form 'x-amz-signature="$(jq '."x-amz-signature"' imageupload.json)"' --form 'Content-MD5="$(jq '."Content-MD5"' imageupload.json)"' --form 'file=@\"./${imagename}.tar\"' -o response_second.json" > response.sh
chmod +x response.sh

cat response.sh
./response.sh

#cat response_second.json

sleep ${timeout}

app_container_uuid=$(jq .container_image_uuid resource/appbundle/appbundle.json)
echo "UUID for app: $app_container_uuid"

img_container_uuid=$(jq .uuid imageupload.json | tr -d '"')
echo "UUID for new image: $img_container_uuid"

eval $(jq '.container_image_uuid="'${img_container_uuid}'"' resource/appbundle/appbundle.json > app_v1.json)
cp app_v1.json resource/appbundle/appbundle.json

app_container_uuid=$(jq .container_image_uuid resource/appbundle/appbundle.json)
echo "UUID for app: $app_container_uuid"