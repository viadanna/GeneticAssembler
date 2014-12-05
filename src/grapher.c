#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int readsequences(char** lines, char* filename, char**aux) {
	int size, i = 0;
	FILE* read;
	char *buffer;

	read = fopen (filename, "r");
	fseek(read, 0, SEEK_END);
	size = ftell(read);
	buffer = malloc (size+1);
	rewind(read);
	fread (buffer, size, 1, read);

	*aux = strtok(buffer, "\n");
	while (*aux != NULL) {
		if ((**aux) == '>') {
			lines[i] = strtok(NULL, "\n");
			i++;
		}
		else {
			memmove (&(lines[i-1][strlen(lines[i-1])]), *aux, strlen(*aux));
		}
		*aux = strtok(NULL, "\n");
	}

	return i;
}

int max (int a, int b, int c) {
	 
	 int m = a;
     (m < b) && (m = b);
     (m < c) && (m = c);
     return m;
}

int sp_align (char* a, char* b) {
	int match = 3, mismatch = -2, gap = -5, i, j, result;
	int m = strlen(a)+1, n = strlen(b)+1;
	int** dp;
	dp = malloc(m*sizeof(int*));
	for (i = 0; i < m; i++) {
		dp[i] = calloc (n, sizeof(int));
	}

	for (j = 1; j < n; j++) {
		dp[0][j] = j * gap;
	}

	for (i = 1; i < m; i++) {
		for (j = 1; j < n; j++) {
			dp[i][j] = max(dp[i-1][j]+gap,dp[i][j-1]+gap,dp[i-1][j-1]+(a[i-1] == b[j-1] ? match : mismatch));
		}
	}

	result = dp[i-1][0];
	for (j = 1; j < n; j++) {
		(result < dp[i-1][j]) && (result = dp[i-1][j]);
	}

	for (i = 0; i < m; i++) {
		free (dp[i]);
	}
	free (dp);

	return result;
}


int main () {
	int i, j, n;
	char** f, *aux;

	f = malloc(sizeof(char*)*1000);

	n = readsequences(f, "reads.fasta", &aux);

	for (i = 0; i < n; i++) {
		for (j = 0; j < n; j++) {
			if (i != j) {
				printf("%d %d %d\n", i+1, j+1, sp_align(f[i], f[j]));
			}
		}
	}

	free (aux);
	free (f);

	return 0;
}