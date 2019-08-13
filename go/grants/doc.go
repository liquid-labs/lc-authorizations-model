// The grants package provides functions for dealing with user grants over system entities. Grants may be given to Users or User Groups, which effectively provide 'roles'. Grants on Containers confer the grant to all contained items.
package grants

// TODO: if we get a good UC, we could use cookies to implement negative grants. The lower-level grant will be found first, and the special 'NEGATIVE' or 'RESCEND' cookie checked by the base library.

// TODO: make sure user or user group is active.
