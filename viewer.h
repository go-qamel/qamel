#pragma once

#ifndef QAMEL_VIEWER_H
#define QAMEL_VIEWER_H

#ifdef __cplusplus

// Class
class QamelView;

extern "C" {
#endif

// Constructor
void* Viewer_NewViewer();

// Methods
void Viewer_SetSource(void* ptr, char* url);
void Viewer_SetResizeMode(void* ptr, int resizeMode);
void Viewer_SetFlags(void* ptr, int flags);
void Viewer_SetHeight(void* ptr, int height);
void Viewer_SetWidth(void* ptr, int width);
void Viewer_SetMaximumHeight(void* ptr, int height);
void Viewer_SetMaximumWidth(void* ptr, int width);
void Viewer_SetMinimumHeight(void* ptr, int height);
void Viewer_SetMinimumWidth(void* ptr, int width);
void Viewer_SetOpacity(void* ptr, double opacity);
void Viewer_SetTitle(void* ptr, char* title);
void Viewer_SetVisible(void* ptr, bool visible);
void Viewer_SetPosition(void* ptr, int x, int y);
void Viewer_SetIcon(void* ptr, char* fileName);
void Viewer_Show(void* ptr);
void Viewer_ShowMaximized(void* ptr);
void Viewer_ShowMinimized(void* ptr);
void Viewer_ShowFullScreen(void* ptr);
void Viewer_ShowNormal(void* ptr);
void Viewer_SetWindowStates(void* ptr, int state);
void Viewer_ClearComponentCache(void* ptr);
void Viewer_Reload(void* ptr);

#ifdef __cplusplus
}
#endif

#endif