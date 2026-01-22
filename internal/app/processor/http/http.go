package rprocessor

// Обратите внимание на название пакета
// Приставка r (m,p) будет добавляться во многих других пакетах.

import (
	"catalog-service/internal/app/config/config"
	"net/http"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth handler.Health, cfg section.ProcessorWebServer) *httpProc {
	// создаем мультиплексор

	// регистрируем HealthCheck

	// здесь будет регистрация остальных хэндлеров
	// TODO: добавить регистрацию продуктов, категорий и т.д.

	// обходим маршруты для дебага через r.Walk
	// если не получится реализовать r.Walk() просто добавляем лог, когда регистрируем маршрут

	// создаем сервер и возвращаем его

	return &s
}
