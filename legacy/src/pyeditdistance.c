/* make functions return Py_ssize_t */
#define PY_SSIZE_T_CLEAN

#include <Python.h>

#include "editdistance.h"

static PyObject *
Distance(PyObject *self, PyObject *args, PyObject *keywds)
{
    char *a;
    char *b;
    Py_ssize_t a_size;
    Py_ssize_t b_size;

    static char *kwlist[] = {"a", "b", NULL};

    if (!PyArg_ParseTupleAndKeywords(args, keywds,
                                     "s#s#",
                                     kwlist,
                                     &a, &a_size,
                                     &b, &b_size)) {
        return NULL;
    }

    printf("%d\n", editdistance(a, (size_t) a_size, b, (size_t) b_size));

    return Py_BuildValue("i", editdistance(a, (size_t) a_size, b, (size_t) b_size));
}

static PyMethodDef EditDistanceMethods[] = {
    {"editdistance", (PyCFunction)Distance,METH_VARARGS | METH_KEYWORDS,
     "Calculate the edit distance between two strings."},
    {NULL, NULL, 0, NULL}
};

PyMODINIT_FUNC
initpyeditdistance(void)
{
    (void) Py_InitModule("pyeditdistance", EditDistanceMethods);
}
