pipeline {
    agent {label 'mp-vm'}

    stages {
        stage('Step 1: Download Git repo') {
          when {
              expression {
                image_upload_required.toBoolean() || icon_upload_required.toBoolean() || lua_upload_required.toBoolean()
              }
          }
            steps {
              sh '''
rm -Rf aruba-iotops-example-ble

git clone --depth 1 --single-branch -b main ${github_repo}
                '''
            }
        }

        stage('Step 2: Build App image and Update ADP App Image') {
          when {
            expression {
              image_upload_required.toBoolean()
            }
          }
            steps {
                sh '''
eval ${build_command}

docker save ${gitimagename}:${gitimageversion} > ${imagename}.tar

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
chmod 777 response.sh

cat response.sh
./response.sh

cat response_second.json

sleep ${timeout}
pwd

app_container_uuid=$(jq .container_image_uuid ${app_bundle_path})
echo "UUID for app: $app_container_uuid"

img_container_uuid=$(jq .uuid imageupload.json | tr -d '"')
echo "UUID for new image: $img_container_uuid"

#sed -i 's/"container_image_uuid":${app_container_uuid}/"container_image_uuid":${img_container_uuid}/g' ${app_bundle_path}
#eval $(sed -i 's/"container_image_uuid":'${app_container_uuid}'/"container_image_uuid":'${img_container_uuid}'/g' ${app_bundle_path})
eval $(jq '.container_image_uuid="'${img_container_uuid}'"' ${app_bundle_path} > app_v1.json)
cp app_v1.json ${app_bundle_path}

app_container_uuid=$(jq .container_image_uuid ${app_bundle_path})
echo "UUID for app: $app_container_uuid"
                '''
            }
        }

        stage('Step 3: Update ADP App icon') {
          when {
            expression {
              icon_upload_required.toBoolean()
            }
          }
            steps {
                sh '''
if [ ! -d ${app_icon_path} ]; then
    app_icon=$(base64 -w 0 ${app_icon_path})
    eval $(jq '.icon.url=null | .icon.data="'${app_icon}'"' ${app_bundle_path} > app_v2.json)
    cp app_v2.json ${app_bundle_path}
fi
                '''
            }
        }

        stage('Step 4: Update ADP App lua script') {
          when {
            expression {
              lua_upload_required.toBoolean()
            }
          }
            steps {
                sh '''
if [ ! -d ${lua_file_path} ]; then
    lua_base=$(base64 -w 0 ${lua_file_path})
    eval $(jq '.lua_script.file_id=null | .lua_script.data="'${lua_base}'"' ${app_bundle_path} > app_v3.json)
    cp app_v3.json ${app_bundle_path}
fi
                '''
            }
        }

        stage('Step 5: Update ADP App') {
          when {
              expression {
                image_upload_required.toBoolean() || icon_upload_required.toBoolean() || lua_upload_required.toBoolean()
              }
          }
            steps {
                sh '''
prefix="'"

eval "curl '${url}/iot_operations/api/v1/adp/apps/${appid}/versions?pageNumber=1&pageSize=1000' \
-H 'authorization: Bearer ${token}' -o appversion.json"

version=$(jq '.content | .[] | select (.status == "DRAFT")'.version appversion.json)

eval "curl '${url}/iot_operations/api/v1/adp/apps/draft/${appid}/version/${version}' \
  -X POST \
  -H 'authorization: Bearer ${token}' \
  -H 'content-type: application/json' \
  -H 'referer: ${url}/swagger/central/' \
  --data-raw ${prefix}$(cat $app_bundle_path)${prefix} \
  -o appupdate.json" > response.sh

cat appupdate.json
                '''
            }
          
        }
    }
}

