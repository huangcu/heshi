# heshi

declare -x http_proxy=xxx
declare -x https_proxy=xx
gb vendor restore
gb build
./bin/heshi_service


http://localhost:8443/api/users
post
Form: (required fields)
cellphone / email
password
user_type


bd21deed-b8b7-424b-b462-1db71bd03dbe

http://localhost:8443/api/users/:id
Get
:id: UUID


http://localhost:8443/api/admin/users
get

get all users
