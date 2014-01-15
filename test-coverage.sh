go list all | grep "cover"

echo "mode: set" > acc.out
for Dir in $(find ./* -maxdepth 10 -type d ); 
do
   go test -coverprofile=profile.out $Dir
   cat profile.out | grep -v "mode: set" >> acc.out
done
goveralls -coverprofile=acc.out $COVERALLS-KEY
rm -rf ./profile.out
rm -rf ./acc.out
