class Solution:
    def lengthOfLongestSubstring(self, s: str) -> int:
        new = {}
        left = 0
        length = 0

        for right, char in enumerate(s):
            if char in new and left <= new[char]:
                left = new[char] + 1
            new[char] = right
            length = max(length, right - left + 1)

        return length