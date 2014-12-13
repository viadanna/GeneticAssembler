/* Program to test the edit distance algorithm */

#include <stdio.h>

#include "../editdistance.h"

int main (void)
{
    char a[] = "genoma";
    char b[] = "gnomos";

    int dist = editdistance(a, sizeof(a)-1, b, sizeof(b)-1);
    if (dist != 4) {
        printf("editdistance test: Failed. (%d != %d)\n", 4, dist);
        return 1;
    }

    puts("editdistance test: Passed.");
    return 0;
}
