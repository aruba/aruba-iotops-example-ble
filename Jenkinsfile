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
rm -Rf ${project_name}

git clone --depth 1 --single-branch -b ${github_local_branch} ${github_repo}
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
                sh './${project_name}/resource/scripts/update_containerimage.sh'
            }
        }

        stage('Step 3: Update ADP App icon') {
          when {
            expression {
              icon_upload_required.toBoolean()
            }
          }
            steps {
                sh './${project_name}/resource/scripts/update_icon.sh'
            }
        }

        stage('Step 4: Update ADP App lua script') {
          when {
            expression {
              lua_upload_required.toBoolean()
            }
          }
            steps {
                sh './${project_name}/resource/scripts/update_appbundle.sh'
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
cd ${project_name}
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

