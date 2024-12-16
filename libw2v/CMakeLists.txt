cmake_minimum_required(VERSION 3.10)
project (w2v)
 
set(PROJECT_ROOT_DIR ${CMAKE_SOURCE_DIR})
set(LIBRARY_OUTPUT_PATH ${PROJECT_ROOT_DIR}/lib)
 
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_FLAGS_DEBUG "-O0 -ggdb")
set(CMAKE_CXX_FLAGS_RELEASE "-Ofast -march=native -funroll-loops -ftree-vectorize -DNDEBUG")

set(PROJECT_INCLUDE_DIR ${PROJECT_ROOT_DIR}/include)

if (${CMAKE_SYSTEM_NAME} MATCHES "Darwin")
    message("-- Detecting Darwin - done")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra -Wpedantic -Werror")
endif()

if (${CMAKE_SYSTEM_NAME} MATCHES "Linux" OR ${CMAKE_SYSTEM_NAME} MATCHES "FreeBSD")
    message("-- Detecting Linux - done")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra -Wpedantic -Werror")
    set(LIBS "-pthread")
endif()

if (${CMAKE_SYSTEM_NAME} MATCHES "Linux")
    set(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} -s")
endif()

add_subdirectory(src)

set(CPACK_SOURCE_GENERATOR "TGZ")
set(CPACK_PACKAGE_VERSION_MAJOR "0")
set(CPACK_PACKAGE_VERSION_MINOR "0")
set(CPACK_PACKAGE_VERSION_PATCH "0")
set(CPACK_SOURCE_IGNORE_FILES
  \\.git/
  build/
  ".*~$"
)
set(CPACK_VERBATIM_VARIABLES YES)
include(CPack)
