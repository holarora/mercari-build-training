from typing import List
class ListNode:
    def __init__(self, value=0, next=None):
        self.value = value
        self.next = next


class Solution(object):
    #O(m+n) time and O(n) space
    def getIntersectionNode(self, headA: ListNode, headB: ListNode) -> ListNode:
        nodes_set=set()

        while headA:
            nodes_set.add(headA)
            headA = headA.next

        while headB:
            if headB in nodes_set:
                return headB
            headB = headB.next
        return None

    #O(m+n) time and O(n) space
    def getListLength(self, head:ListNode) -> int:
        length = 0
        while head:
            head = head.next
            length +=1
        return length
    def getIntersectionNode2(self, headA: ListNode, headB: ListNode) -> ListNode:
        lenA, lenB = self.getListLength(headA), self.getListLength(headB)

        while lenA > lenB:
            headA = headA.next
            lenA -= 1
        while lenB > lenA:
            headB = headB.next
            lenB -= 1
        while headA != headB:
            headA = headA.next
            headB = headB.next
        return headA

