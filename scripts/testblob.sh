
echo "\n---head should fail-----"
curl -v -I localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d

echo "\n---get should fail-----"
curl localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d

echo "\n---put----"
curl -X PUT -F "filename=@testblob-e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d" localhost:8080/oci/v1/second/second/blobs/upload

echo "\n---head should ok-----"
curl -v -I localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d

echo "\n---get should ok-----"
curl localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d


echo "\n---delete-----"
curl -X DELETE localhost:8080/oci/v1/second/second/blobs/e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d
