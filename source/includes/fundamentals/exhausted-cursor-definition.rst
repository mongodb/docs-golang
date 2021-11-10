.. note:: 

   A cursor is exhausted when the ``MongoClient`` iterates through its
   results and reaches its last element. Afterwards, the cursor won't
   respond to functions :ref:`to access its elements <cursor-api-docs-golang>`.
