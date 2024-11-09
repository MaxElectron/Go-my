//go:build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) []string {

	// make variable for answer
	ans := make([]string, 0)

	// collect all courses from prereqs into a set
	coursesWithNoPrereqs := make(map[string]bool)
	for _, coursePrereqs := range prereqs {
		for _, course := range coursePrereqs {
			coursesWithNoPrereqs[course] = true
		}
	}

	// remove those which have prereqs
	for course, prereq := range prereqs {
		if len(prereq) != 0 {
			delete(coursesWithNoPrereqs, course)
		} else {
			delete(prereqs, course)
		}
	}

	// make a graph from prereqs
	graph := make(map[string][]string)
	for course, prereq := range prereqs {
		graph[course] = prereq
	}

	// do topological sort by Kahns algorithm

	// make a queue from the courses with no prereqs
	queue := make([]string, 0)
	for course := range coursesWithNoPrereqs {
		queue = append(queue, course)
	}

	// while queue is not empty
	for len(queue) > 0 {
		// remove the course from the queue
		removedCourse := queue[0]
		queue = queue[1:]

		// add the course to the answer
		ans = append(ans, removedCourse)

		// remove the course from all prereqs
		for course, prereqs := range graph {
			for i, prereq := range prereqs {
				if prereq == removedCourse {
					graph[course] = append(prereqs[:i], prereqs[i+1:]...)
					break
				}
			}
		}

		// add courses with no prereqs to the queue
		for course, prereqs := range graph {
			if len(prereqs) == 0 {
				queue = append(queue, course)
				delete(graph, course)
			}
		}
	}

	// if the graph is not empty panic
	if len(graph) > 0 {
		panic("cycle detected")
	}

	return ans
}
