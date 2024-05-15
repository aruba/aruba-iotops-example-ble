pipeline {
    agent any
    parameters {
      string(name: 'url', description: 'Central URL')
      string(name: 'token', description: 'Central NBAPI token')
      string(name: 'appid', description: 'ADP application id')
      booleanParam(name: 'icon_upload_required', description: 'Select the checkbox if app icon update is required for ADP app')
      booleanParam(name: 'lua_upload_required', description: 'Select the checkbox if lua file update is required for ADP app')
      booleanParam(name: 'image_upload_required', description: 'Select the checkbox if container image upload is required for ADP app')
      string(name: 'imagename', defaultValue: 'ExampleApp', description: 'Image name to be uploaded')
      string(name: 'timeout', defaultValue: '120', description: 'Wait time in seconds required during image upload')
    }

    stages {
        stage('Step 1: Build App image and Update ADP App Image') {
          when {
            expression {
              image_upload_required.toBoolean()
            }
          }
            steps {
                sh './resource/scripts/update_containerimage.sh'
            }
        }

        stage('Step 2: Update ADP App icon') {
          when {
            expression {
              icon_upload_required.toBoolean()
            }
          }
            steps {
                sh './resource/scripts/update_icon.sh'
            }
        }

        stage('Step 3: Update ADP App lua script') {
          when {
            expression {
              lua_upload_required.toBoolean()
            }
          }
            steps {
                sh './resource/scripts/update_containerimage.sh'
            }
        }

        stage('Step 4: Update ADP App') {
          when {
              expression {
                image_upload_required.toBoolean() || icon_upload_required.toBoolean() || lua_upload_required.toBoolean()
              }
          }
            steps {
                sh '''
cat resource/appbundle/appbundle.json
prefix="'"

eval "curl '${url}/iot_operations/api/v1/adp/apps/${appid}/versions?pageNumber=1&pageSize=1000' \
-H 'authorization: Bearer ${token}' -o appversion.json"

version=$(jq '.content | .[] | select (.status == "DRAFT")'.version appversion.json)

eval "curl '${url}/iot_operations/api/v1/adp/apps/draft/${appid}/version/${version}' \
  -X POST \
  -H 'authorization: Bearer ${token}' \
  -H 'content-type: application/json' \
  -H 'referer: ${url}/swagger/central/' \
  --data-raw ${prefix}$(cat resource/appbundle/appbundle.json)${prefix} \
  -o appupdate.json" > response.sh

cat appupdate.json
                '''
            }
          
        }
    }
}

