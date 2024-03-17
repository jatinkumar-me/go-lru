import { useState } from "react";

interface CacheData {
  key: string;
  value: string;
}

const BASE_URL = import.meta.env.VITE_BASE_API_URL as string;

function CacheOperations() {
  const [cacheKey, setCacheKey] = useState<number>(0);
  const [cacheValue, setCacheValue] = useState<string>('');

  let canGet: boolean = !isNaN(cacheKey)
  let canSet: boolean = canGet && cacheValue.length > 0;

  const buttonText = canSet ? "Set cache" : "Get cache";

  const handleGetKey = async () => {
    try {
      if (!BASE_URL) {
        console.error("Base url not defined");
        return;
      }
      const response = await fetch(`${BASE_URL}/cache/get?key=${cacheKey}`);
      if (response.status != 200) {
        setCacheValue("");
        return;
      }
      const data: CacheData = await response.json();
      setCacheValue(data.value);
    } catch (error) {
      console.error('Error fetching key:', error);
    }
  };

  const handleSetKeyValue = async () => {
    if (!cacheValue.trim()) {
      alert('Please enter a value');
      return;
    }

    try {
      if (!BASE_URL) {
        console.error("Base url not defined");
        return;
      }
      const response = await fetch(`${BASE_URL}/cache/set`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ key: cacheKey, value: cacheValue }),
      });
      if (response.ok) {
        console.log('Key/value pair set successfully');
        setCacheValue("");
      } else {
        console.error('Failed to set key/value pair');
        setCacheValue("");
      }
    } catch (error) {
      console.error('Error setting key:', error);
    }
  };

  function handleClick(): void {
    if (canSet) {
      handleSetKeyValue();
    } else {
      handleGetKey();
    }
  }

  return (
    <div>
      <h1>Cache Operations</h1>
      <div>
        <fieldset>
          <legend><h2>Set Key/Value in Cache</h2></legend>
          <div>
            <label
              htmlFor="cache-key"
            >
              Cache Key
            </label>
            <input
              type="number"
              id="cache-key"
              value={cacheKey}
              onChange={(e) => setCacheKey(parseInt(e.target.value))}
              placeholder="Enter key"
            />

            <label
              htmlFor="cache-value"
            >
              Cache Value
            </label>
            <input
              type="text"
              id="cache-value"
              value={cacheValue}
              onChange={(e) => setCacheValue(e.target.value)}
              placeholder="Enter value"
            />
          </div>
          <button onClick={handleClick}>{buttonText}</button>
        </fieldset>
      </div>
    </div>
  );
}

export default CacheOperations;
