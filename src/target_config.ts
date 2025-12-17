export const target_tauri = false;

export const api_proxy_addr = "http://10.111.128.121:8080";

export const dest_api = (target_tauri) ? api_proxy_addr + "/api" : "/api";

export const dest_root = (target_tauri) ? "" : "/acid-app";