class Solution(object):
    def wordPattern(self, pattern: str, s: str) -> bool:
        pattern1 = ""
        words = s.split()
        word_dict = {}
        word_inverted_dict = {}

        if len(words) != len(pattern):
            return False
        for i, p in enumerate(pattern):
            if p not in word_dict and words[i] not in word_inverted_dict:
                word_dict[p] = words[i]
                word_inverted_dict[words[i]] = p
            elif p in word_dict and words[i] in word_inverted_dict:
                if word_dict[p] != words[i]:
                    return false
            else:
                return false
        return True

if __name__ == "__main__":
    pattern = "abba"
    s = "dog cat cat dog"
    result = Solution().wordPattern(pattern, s)
    print(result)

