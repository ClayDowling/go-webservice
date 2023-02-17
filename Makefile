.PHONY: test package deploy acceptance-test

webservice: test
	go build

test:
	go test

package: acceptance-test
	zip webservice.zip webservice

deploy: acceptance-test
	scp webservice 

acceptance-test: webservice
	rm -f users.db
	sqlite3 users.db < data/users.sql
	./webservice &
	sleep 2
	cd acceptance && go test 
	WSPID=$$(<webservice.pid) && kill $$WSPID
