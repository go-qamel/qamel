#pragma once

#ifndef QAMEL_ENGINE_H
#define QAMEL_ENGINE_H

#ifdef __cplusplus
extern "C" {
#endif

// Constructors
void* Engine_NewEngine();
void* Engine_NewEngineWithSource(char* source);

// Methods
void Engine_Load(void* ptr, char* url);
void Engine_ClearComponentCache(void* ptr);

#ifdef __cplusplus
}
#endif

#endif