http GET ':5005/api/posts?lat=41.947259&lon=-87.65438' 'Authorization:Bearer 1FETidvRvKNBS7oBQFrLnSMlBeX'

http POST :5005/api/post 'Authorization:Bearer 1FETidvRvKNBS7oBQFrLnSMlBeX'  lat:=41.95 lon:=-87.66 accessType=public 'content={"title": "First", "body": "New phone who dis"}'

http POST :5005/api/post/1FEWViwSeKkQ8hqaVkM2crOezbj/vote?state=1 'Authorization:Bearer 1FETidvRvKNBS7oBQFrLnSMlBeX'