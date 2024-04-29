cd ${project_name}
if [ -d "lua" ]; then
    lua_base=$(base64 -w 0 lua/ibeacon.lua)
    eval $(jq '.lua_script.file_id=null | .lua_script.data="'${lua_base}'"' resource/appbundle/appbundle.json > app_v3.json)
    cp app_v3.json resource/appbundle/appbundle.json
fi