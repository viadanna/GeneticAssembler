#include <stdlib.h>
#include <stdio.h>

#include "editdistance.h"

static const int match = 0;
static const int gap = -1;
static const int mismatch = -2;

/* cost of aligning a with b */
static int align (char a, char b)
{
    return ((a == b) ? match : mismatch);
}

static char max (char a, char b, char c)
{
   char max = ( a < b ) ? b : a;
   return ( ( max < c ) ? c : max );
}

int editdistance (char* a, size_t m, char* b, size_t n)
{
    int i, j;
    int result;
    /* a matrix allocated in linear form. We use M[i * m + j] to access line i
       and column j */
    m++;
    n++;
    int *M = malloc(m*n*sizeof(int));

    /* initialize */
    for (i = 0; i < n; i++) {
        M[i * m] = gap * i;
    }

    for (j = 0; j < m; j++) {
        M[j] = gap * j;
    }

    /* calculate global alignment */
    for (i = 1; i < n; i++) {
        for (j = 1; j < m; j++) {
            M[i * m + j] = max(M[i * m + j-1] + gap,
                               M[(i-1) * m + j-1] + align(a[j-1], b[i-1]),
                               M[(i-1) * m +j] + gap);
        }
    }

    result = -1 * M[(n-1) * m + (m-1)];
#ifdef DEBUG
    for (i = 0; i < n; i++) {
        for (j = 0; j < m; j++) {
            printf("%d\t", M[i * m + j]);
        }
        printf("\n");
}
#endif

    free(M);

    return result;
}
