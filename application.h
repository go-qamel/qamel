#pragma once

#ifndef QAMEL_APPLICATION_H
#define QAMEL_APPLICATION_H

#ifdef __cplusplus
extern "C" {
#endif

// Constructor
void* App_NewApplication(int argc, char* argv[]);

// Static function
void App_SetAttribute(int attribute, bool on);

// Method
void App_SetFont(void *ptr, char *family, int pointSize, int weight, bool italic);
void App_SetQuitOnLastWindowClosed(void* ptr, bool quit);
void App_SetApplicationDisplayName(void* ptr, char* name);
void App_SetWindowIcon(void* ptr, char* fileName);
void App_SetApplicationName(void* ptr, char* name);
void App_SetApplicationVersion(void* ptr, char* version);
void App_SetOrganizationName(void* ptr, char* name);
void App_SetOrganizationDomain(void* ptr, char* domain);
int App_Exec(void* ptr);

#ifdef __cplusplus
}
#endif

#endif