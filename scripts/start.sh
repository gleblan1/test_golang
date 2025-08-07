#!/bin/bash

# Скрипт для запуска Crypto Price Tracker

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Функция для вывода сообщений
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
}

# Проверяем наличие Docker
if ! command -v docker &> /dev/null; then
    error "Docker не установлен. Пожалуйста, установите Docker."
    exit 1
fi

# Проверяем наличие Docker Compose
if ! command -v docker-compose &> /dev/null; then
    error "Docker Compose не установлен. Пожалуйста, установите Docker Compose."
    exit 1
fi

# Функция для остановки сервисов
cleanup() {
    log "Останавливаем сервисы..."
    docker-compose down
    exit 0
}

# Обработка сигналов для graceful shutdown
trap cleanup SIGINT SIGTERM

# Основная логика
main() {
    log "Запуск Crypto Price Tracker..."
    
    # Создаем .env файл если его нет
    if [ ! -f .env ]; then
        log "Создаем .env файл из примера..."
        cp env.example .env
    fi
    
    # Запускаем сервисы
    log "Запускаем Docker Compose сервисы..."
    docker-compose up -d postgres
    
    # Ждем готовности базы данных
    log "Ожидаем готовности базы данных..."
    sleep 10
    
    # Применяем миграции
    log "Применяем миграции базы данных..."
    docker-compose run --rm migrate
    
    # Запускаем API и Worker
    log "Запускаем API и Worker сервисы..."
    docker-compose up -d api worker
    
    log "Crypto Price Tracker запущен!"
    log "API доступен по адресу: http://localhost:8080"
    log "Swagger UI: http://localhost:8080/swagger/"
    log "Health check: http://localhost:8080/health"
    log ""
    log "Для просмотра логов выполните: docker-compose logs -f"
    log "Для остановки выполните: docker-compose down"
    
    # Показываем логи
    docker-compose logs -f
}

# Запуск основной функции
main 