@echo off
setlocal enabledelayedexpansion

:: Початкові координати (від 0 до 100, щоб уникнути проблем з плаваючою комою)
set /a x=50
set /a y=50

echo Starting animation loop...
echo Press Ctrl+C to stop.

:loop
    :: Конвертуємо цілі числа в рядок формату 0.xx
    set "x_float=0.!x!"
    set "y_float=0.!y!"

    echo Drawing figure at !x_float!, !y_float!

    :: Динамічно створюємо файл з командами для поточного кадру.
    :: Це найнадійніший спосіб уникнути проблем з форматуванням.
    (
        echo white
        echo figure !x_float! !y_float!
        echo update
    ) > scripts\current_frame.txt

    :: Відправляємо створений файл на сервер
    curl -s -X POST --data-binary @scripts/current_frame.txt http://localhost:17000

    :: Готуємо координати для наступного кадру
    set /a x=!x!+2
    set /a y=!y!+2

    :: Якщо фігура виходить за межі, починаємо спочатку
    if !x! gtr 98 set x=2
    if !y! gtr 98 set y=2
    
    :: Пауза в 1 секунду
    timeout /t 1 /nobreak >nul
    
    :: Повторюємо цикл
    goto loop