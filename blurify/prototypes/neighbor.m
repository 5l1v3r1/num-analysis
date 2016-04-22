function [m] = neighbor(d)
  m = eye(d^2, d^2);
  for i = 1:d
    for j = 1:d
      rowIdx = pointIdx(d, i, j);
      m(rowIdx, rowIdx) = 1/5;
      m(rowIdx, pointIdx(d, i-1, j)) = 1/5;
      m(rowIdx, pointIdx(d, i+1, j)) = 1/5;
      m(rowIdx, pointIdx(d, i, j-1)) = 1/5;
      m(rowIdx, pointIdx(d, i, j+1)) = 1/5;
    end
  end
end

function [idx] = pointIdx(d, x, y)
  x = wrapIdx(d, x);
  y = wrapIdx(d, y);
  idx = x + (y-1)*d;
end

function [idx] = wrapIdx(d, x)
  idx = x;
  if (idx < 1)
    idx = d + idx;
  elseif (idx > d)
    idx = idx - d;
  end
end
