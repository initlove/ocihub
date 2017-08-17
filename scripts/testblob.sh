# the session id should be replaced by using 'post' operation, we can get 'Session-Id' in the header
# call ./create_session_id.sh
SESSION_ID=5b3188e3-6eaa-4fc9-9e84-dd43a62682fb

echo "\n---head should fail-----"
curl -I localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d

echo "\n---get should fail-----"
curl localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d

echo "\n---post----"
curl -X POST localhost:8080/oci/v1/second/second/blobs/uploads/

echo "\n---patch---"
curl -F "filename=@testblob-e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d" -X PATCH localhost:8080/oci/v1/second/second/blobs/uploads/$SESSION_ID

echo "\n---put----"
curl -X PUT localhost:8080/oci/v1/second/second/blobs/uploads/$SESSION_ID

echo "\n---put wil fail since the session id is released----"
curl -X PUT localhost:8080/oci/v1/second/second/blobs/uploads/$SESSION_ID

echo "\n---head should ok-----"
curl -I localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d

echo "\n---get should ok-----"
curl localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d


echo "\n---delete-----"
curl -X DELETE localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d
