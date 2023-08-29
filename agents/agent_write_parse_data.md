El siguiente mensaje es un listado de datos climáticos del proyecto Redmop. Redmop es un sistema de monitoreo climático distribuido que utiliza una red de voluntarios para recopilar datos sobre las condiciones meteorológicas en una zona geográfica específica. Los datos que el mensaje debe contener son:

- Código de la estación (ej.: RC-Cy-CP-34)
- Quebrada y distrito (ej.: CULTURA Y PROGRESO CHACLACAYO)
- Intensidad de la lluvia (en una escala de 1 a 5, ej.: 3)
- Nivel de lluvia acumulada (en mm, ej.: 0.5mm)
- Tiempo actual (ej.: Cielo nublado sigue llovizna)
- Fecha y hora de la revisión del pluviómetro (ej.: 09-06-2023 7:31am)

Si no se puede parsear la información, responde solo el string "Hubo un error al parsear la data".

No añadas explicación, solo los datos.
La respuesta debe ser en formato JSON, con los siguientes campos:

- "station_code": string,
- "stream_name": string,
- "rain_intensity": int,
- "rain_level": float,
- "current_weather": string,
- "date_time": string (ISO format)

Si se puede parsear, pero tiene campos incompletos, genera un JSON con los campos faltantes como "null"
