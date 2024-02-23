package api

type shortBody struct {
	LongURL string `json:"long_url"`
}


/*
	Shortener logic

	We receive https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74

	1- convert this with a hash function so each url.String() is mapped with a short hash value

	2- store the mapping in the database

	3- return the shorten URL


	Redirect logic

	We receive https://<my-domain>/<slug>

	1- inspect the database for such a slug value

		if found, return the associated url value and send an http 301 response with the url in location

		if not found return 404



*/



func shorten(id int) (string, error) {
	
	/* a goode idea for the shortener is to generate an unique ID associated to the URL (from the DB)
	
		And then return this ID encoded in another base like base62 to get a shorten version of it

		The downside is the slug size will grow constantly as new URL will be encoded
		(base62 is shorter than base64)
	*/

	

	return " ",  nil 

}


func 
