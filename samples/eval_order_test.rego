# This policy asserts that the rules are not evaluated in parallel (that is, print statements are executed in order)
package example

import future.keywords.if

allow if {
	# time-consuming statement
	http.send({"method": "get", "url": "https://www.google.com"})

	print("✅allow 1, this should be printed")
}

allow if {
	print("❌allow 2, this should NOT be printed")
	false
}

allow if {
	print("❌allow 3, this should NOT be printed, of course")
	false
}
