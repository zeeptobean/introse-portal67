package main

import (
	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/service"
)

func main() { service.Main(routing.Academic, "8082") }
