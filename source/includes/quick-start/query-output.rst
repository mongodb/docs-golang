When you run ``main.go``, it should output the details of the
movie from the sample dataset which will look something like this:

.. code-block:: json

   [
       {
           "Key": "_id",
           "Value": "573a1398f29313caabce9682"
       },
       {
           "Key": "plot",
           "Value": "A young man is accidentally sent 30 years into the past in a time-traveling DeLorean invented by his friend, Dr. Emmett Brown, and must make sure his high-school-age parents unite in order to save his own existence."
       },
       ...
       {
           "Key": "title",
           "Value": "Back to the Future"
       },
       ...
   ]

If you receive no output or an error, check whether you included the proper
connection string in your ``main.go`` file, and whether you loaded the
sample dataset in your MongoDB Atlas cluster
