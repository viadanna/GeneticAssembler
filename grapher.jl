"
    Grapher: Generate the edges of Overlap Graph from given reads
    Author: Paulo Viadanna
"

readsequences(source::String) = [join(split(i, '\n')[2:end]) for i in split(readall(source), ">")]

" Suffix-Prefix alignment based on Newman-Munsch "
function sp_align(a, b, match=3, mismatch=-2, gap=-5)
    # Array and acessory variables initialization
    m, n = endof(a)+1, endof(b)+1
    dp = zeros(Int32, m, n)

    # Initialize dynamic programming array
    for j = 2:n
        dp[1,j] = (j-1) * gap
    end

    # Fill the array
    for i = 2:m
        for j = 2:n
            dp[i,j] = max(dp[i-1,j]+gap,dp[i,j-1]+gap,dp[i-1,j-1]+(a[i-1] == b[j-1] ? match : mismatch))
        end
    end

    # Return the Suffix-Prefix maximum
    return maximum(dp[end,2:end])
end

# Read all sequences into memory
F = readsequences("reads.fasta")

# Loop all sequences and find highest weigth edges for the Overlap Graph
for i in 2:endof(F)
    for j in 2:endof(F)
        if i == j
            continue
        end
        @printf("%d %d %d\n", i, j, sp_align(F[i], F[j]))
    end
end
