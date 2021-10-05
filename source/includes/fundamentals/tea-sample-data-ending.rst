.. tip:: Non-existent Database and Collections

   The driver automatically creates the necessary database and/or collection
   when you perform a write operation against them if they don't already exist.

Each document contains a rating for a type of tea, which corresponds to
the ``type`` and ``rating`` fields.

.. note::

   Each example truncates the ``ObjectID`` value since the driver
   generates them uniquely.
   