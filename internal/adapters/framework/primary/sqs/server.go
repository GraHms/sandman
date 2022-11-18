package sqs

func (gina *Adapter) RUN() {

	err := gina.router.Run(gina.addr)
	if err != nil {
		println("Could not start server", err)
	}
}
