import requests
import time
import random
from datetime import datetime

# Configuration
TICKER = "t"
BASE_URL = "http://localhost:5001" # Adjust port if needed
DURATION_HOURS = 2
INTERVAL_SEC = 1

def simulate():
    print(f"Starting simulation for {TICKER} at {BASE_URL}")
    print(f"Duration: {DURATION_HOURS} hours, Interval: {INTERVAL_SEC}s")
    
    # Get current price to start with
    try:
        res = requests.get(f"{BASE_URL}/api/state/{TICKER.lower()}")
        if res.ok:
            data = res.json()
            last_close = data.get("close") or 24000.0
            last_open = data.get("open") or last_close
            last_high = data.get("high") or last_close
            last_low = data.get("low") or last_close
        else:
            last_close = 24000.0
            last_open = 24000.0
            last_high = 24000.0
            last_low = 24000.0
    except:
        last_close = 24000.0
        last_open = 24000.0
        last_high = 24000.0
        last_low = 24000.0

    start_time = time.time()
    end_time = start_time + (DURATION_HOURS * 3600)
    
    ticks_sent = 0
    try:
        while time.time() < end_time:
            # Simple random walk
            change = random.uniform(-2.0, 2.0)
            last_close += change
            last_high = max(last_high, last_close)
            last_low = min(last_low, last_close)
            
            payload = {
                "date": datetime.now().strftime("%d-%b-%Y"),
                "close": round(last_close, 2),
                "open": round(last_open, 2),
                "high": round(last_high, 2),
                "low": round(last_low, 2),
                "timestamp": int(time.time())
            }
            
            try:
                res = requests.put(f"{BASE_URL}/api/update/{TICKER.lower()}", json=payload)
                if not res.ok:
                    print(f"Error: {res.status_code} - {res.text}")
                else:
                    ticks_sent += 1
                    if ticks_sent % 60 == 0:
                        print(f"Sent {ticks_sent} ticks. Current price: {payload['close']}")
            except Exception as e:
                print(f"Request failed: {e}")
            
            time.sleep(INTERVAL_SEC)
            
    except KeyboardInterrupt:
        print("\nSimulation stopped by user.")

    print(f"Simulation finished. Sent {ticks_sent} ticks.")

if __name__ == "__main__":
    simulate()
