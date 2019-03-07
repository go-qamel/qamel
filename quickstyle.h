#pragma once

#ifndef QAMEL_QUICKSTYLE_H
#define QAMEL_QUICKSTYLE_H

#ifdef __cplusplus
extern "C" {
#endif

void SetQuickStyle(char* style);
void SetQuickStyleFallback(char* style);
void AddQuickStylePath(char* style);

#ifdef __cplusplus
}
#endif

#endif
