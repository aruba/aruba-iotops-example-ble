if [ -d "resource/appicon" ]; then
    pwd
    app_icon=$(base64 -w 0 resource/appicon/appicon.png)
    eval $(jq '.icon.url=null | .icon.data="'${app_icon}'"' resource/appbundle/appbundle.json > app_v2.json)
    cp app_v2.json resource/appbundle/appbundle.json
fi