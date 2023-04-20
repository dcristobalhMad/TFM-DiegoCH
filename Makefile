# Pulumi commands
check:
	cd Infrastructure && pulumi preview

deploy:
	cd Infrastructure && pulumi up

destroy:
	cd Infrastructure && pulumi destroy