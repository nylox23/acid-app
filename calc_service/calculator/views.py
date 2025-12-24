import json
import time
import random
import threading
import requests
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt

INTERNAL_TOKEN = "12345678"
GO_CALLBACK_URL = "http://localhost:8080/api/internal/update-results"

def calculate_logic(data):

    delay = random.randint(5, 10)
    time.sleep(delay)

    results = []

    items = data.get('items', [])
    
    for item in items:
        
        c_mass = float(item.get('c_mass', 0))
        a_mass = float(item.get('a_mass', 0))
        h_plus = float(item.get('h_plus', 1))
        molar = float(item.get('molar', 1))

        if molar == 0:
            val = 0
        else:
            mol_caco3 = c_mass / 100.0
            mol_acid_equiv = (a_mass * h_plus) / (2.0 * molar)
            
            val = 22.4 * min(mol_caco3, mol_acid_equiv)

        results.append({
            "id": item.get('id'),
            "result": val
        })

    try:
        payload = {"results": results}
        headers = {
            "Content-Type": "application/json",
            "X-Internal-Token": INTERNAL_TOKEN
        }
        requests.put(GO_CALLBACK_URL, json=payload, headers=headers)
        print(f"Calculations sent back for {len(results)} items after {delay}s")
    except Exception as e:
        print(f"Error sending callback: {e}")

@csrf_exempt
def calculate_endpoint(request):
    if request.method == 'POST':
        try:
            body = json.loads(request.body)

            thread = threading.Thread(target=calculate_logic, args=(body,))
            thread.start()

            return JsonResponse({"status": "accepted", "message": "Calculation started"}, status=202)
        except Exception as e:
            return JsonResponse({"error": str(e)}, status=400)
    
    return JsonResponse({"error": "Method not allowed"}, status=405)