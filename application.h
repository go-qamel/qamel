#pragma once

#ifndef QAMEL_APPLICATION_H
#define QAMEL_APPLICATION_H

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Constructor
void* App_NewApplication(int argc, char* argv);

// Static function
void App_SetAttribute(long long attribute, bool on);

// Method
void App_SetFont(char *family, int pointSize, int weight, bool italic);
void App_SetQuitOnLastWindowClosed(bool quit);
void App_SetApplicationDisplayName(char* name);
void App_SetWindowIcon(char* fileName);
void App_SetApplicationName(char* name);
void App_SetApplicationVersion(char* version);
void App_SetOrganizationName(char* name);
void App_SetOrganizationDomain(char* domain);
int App_Exec();

#ifdef __cplusplus
}
#endif

#endif